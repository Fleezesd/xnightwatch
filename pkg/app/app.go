package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/component-base/cli"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"

	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
	"github.com/fleezesd/xnightwatch/pkg/version"
)

type App struct {
	name        string
	shortDesc   string
	description string
	run         RunFunc
	cmd         *cobra.Command
	args        cobra.PositionalArgs

	healthCheckFunc HealthCheckFunc

	options CliOptions

	silence bool

	noConfig bool

	// watching and re-loading config files
	watch bool
}

// RunFunc defines the application startup callback function.
type RunFunc func() error

type HealthCheckFunc func() error

type Option func(*App)

func WithOptions(opts CliOptions) Option {
	return func(app *App) {
		app.options = opts
	}
}

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.run = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

// WithHealthCheckFunc is used to set the health check function for the application.
// The app framework will use the function to start a health check server.
func WithHealthCheckFunc(fn HealthCheckFunc) Option {
	return func(app *App) {
		app.healthCheckFunc = fn
	}
}

// WithDefaultHealthCheckFunc set the default health check function.
func WithDefaultHealthCheckFunc() Option {
	fn := func() HealthCheckFunc {
		return func() error {
			go genericoptions.NewHealthOptions().ServeHealthCheck()
			return nil
		}
	}

	return WithHealthCheckFunc(fn())
}

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence() Option {
	return func(app *App) {
		app.silence = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig() Option {
	return func(app *App) {
		app.noConfig = true
	}
}

// WithValidArgs set the validation function to valid non-flag arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// WithWatchConfig watching and re-reading config files.
func WithWatchConfig() Option {
	return func(app *App) {
		app.watch = true
	}
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp(name string, shortDesc string, opts ...Option) *App {
	app := &App{
		name:      name,
		run:       func() error { return nil },
		shortDesc: shortDesc,
	}

	for _, o := range opts {
		o(app)
	}

	return app
}

// buildCommand is used to build a cobra command.
func (app *App) buildCommand() {
	cmd := &cobra.Command{
		Use:   formatBaseName(app.name),
		Short: app.shortDesc,
		Long:  app.description,
		RunE:  app.runCommand,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Args: app.args,
	}

	if !cmd.SilenceUsage {
		cmd.SilenceUsage = true
		// flag error handling
		cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
			c.SilenceUsage = false
			return err
		})
	}
	// In all cases error printing is done below
	cmd.SilenceErrors = true

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true

	var fss cliflag.NamedFlagSets
	if app.options != nil {
		fss = app.options.Flags()
	}

	version.AddFlags(fss.FlagSet("global"))
	if !app.noConfig {
		AddConfigFlag(fss.FlagSet("global"), app.name, app.watch)
	}

	for _, f := range fss.FlagSets {
		cmd.Flags().AddFlagSet(f)
	}

	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmd, fss, cols)

	app.cmd = cmd
}

// formatBaseName is formatted as an executable file name under different
// operating systems according to the given name.
func formatBaseName(name string) string {
	if runtime.GOOS == "windows" {
		name = strings.ToLower(name)
		name = strings.TrimSuffix(name, ".exe")
	}
	return name
}

// Run is used to launch the application.
func (app *App) Run() {
	os.Exit(cli.Run(app.cmd))
}

func (app *App) runCommand(cmd *cobra.Command, args []string) error {
	// display application version information
	version.PrintAndExitIfRequested(app.name)

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if app.options != nil {
		if err := viper.Unmarshal(app.options); err != nil {
			return err
		}

		// set default options
		if err := app.options.Complete(); err != nil {
			return err
		}

		// validate options
		if err := app.options.Validate(); err != nil {
			return err
		}
	}

	// init log
	log.Init(logOptions())
	defer log.Sync() // Sync flushes any buffered log entries.

	if !app.silence {
		log.Infow("Starting application", "name", app.name, "version", version.Get().ToJSON())
		log.Infow("Golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))

		if !app.noConfig {
			PrintConfig()
		} else if app.options != nil {
			cliflag.PrintFlags(cmd.Flags())
		}
	}

	if app.healthCheckFunc != nil {
		if err := app.healthCheckFunc(); err != nil {
			return err
		}
	}

	// run application
	return app.run()
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		EnableColor:       viper.GetBool("log.enable-color"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

func init() {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "console")
	viper.SetDefault("log.output-paths", []string{"stdout"})
}
