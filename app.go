package glue

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var server http.Server

func Handle(c *Config) error {
	keyGenerate := flag.NewFlagSet("key:generate", flag.ExitOnError)
	serve := flag.NewFlagSet("serve", flag.ExitOnError)

	if !flag.Parsed() {
		flag.Parse()
	}
	args := flag.Args()
	// first element is the name of the command
	if len(args) > 0 {
		args = args[1:]
	}

	switch flag.Arg(0) {
	case keyGenerate.Name():
		show := keyGenerate.Bool("show", false, "Display the key instead of modifying files.")
		if err := keyGenerate.Parse(args); err != nil {
			return err
		}
		key, err := generateRandomKey()
		if err != nil {
			return err
		}

		if *show {
			fmt.Println(key)
			return nil
		}

		if currentKey := c.Key; currentKey != "" {
			return nil
		}
		if err := writeNewEnvironmentFileWith(key); err != nil {
			return err
		}

		fmt.Println("Application key set successfully.")
	case serve.Name():
		addr := serve.String("addr", ":8080", `Specify the TCP address for the server to listen on, in the form "host:port".`)
		if err := serve.Parse(args); err != nil {
			return err
		}

		auth := http.NewServeMux()
		auth.HandleFunc("/login", func(resp http.ResponseWriter, req *http.Request) {
			_, _ = fmt.Fprintln(resp, "login")
		})
		auth.HandleFunc("/logout", func(resp http.ResponseWriter, req *http.Request) {
			_, _ = fmt.Fprintln(resp, "logout")
		})
		if c.Auth.Register {
			auth.HandleFunc("/register", func(resp http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprintln(resp, "register")
			})
		}
		if c.Auth.Reset {
			auth.HandleFunc("/password/reset", func(resp http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprintln(resp, "password/reset")
			})
			auth.HandleFunc("/password/email", func(resp http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprintln(resp, "password/email")
			})
		}
		if c.Auth.Confirm {
			auth.HandleFunc("/password/confirm", func(resp http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprintln(resp, "password/confirm")
			})
		}
		if c.Auth.Verify {
			auth.HandleFunc("/email/verify", func(resp http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprintln(resp, "email/verify")
			})
			auth.HandleFunc("/email/resend", func(resp http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprintln(resp, "email/resend")
			})
		}

		server = http.Server{Addr: *addr, Handler: auth}

		if err := server.ListenAndServe(); err != nil {
			return err
		}
	}

	return nil
}

type AuthConfig struct {
	Register bool
	Reset    bool
	Confirm  bool
	Verify   bool
}

type Config struct {
	Key  string
	Auth AuthConfig
}

func NewConfig() *Config {
	return &Config{
		Key: os.Getenv("APP_KEY"),
		Auth: AuthConfig{
			Register: true,
			Reset:    true,
		},
	}
}

func Shutdown(ctx context.Context) error {
	return server.Shutdown(ctx)
}
