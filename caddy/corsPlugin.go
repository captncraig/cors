package caddy

import (
	"github.com/captncraig/cors"
	"net/http"
	"strconv"
	"strings"

	"github.com/mholt/caddy/config/setup"
	"github.com/mholt/caddy/middleware"
)

type corsRule struct {
	Conf *cors.Config
	Path string
}

func Setup(c *setup.Controller) (middleware.Middleware, error) {
	rules, err := parseRules(c)
	if err != nil {
		return nil, err
	}
	return func(next middleware.Handler) middleware.Handler {
		return middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
			for _, rule := range rules {
				if middleware.Path(r.URL.Path).Matches(rule.Path) {
					rule.Conf.HandleRequest(w, r)
					if cors.IsPreflight(r) {
						return 200, nil
					}
					break
				}
			}
			return next.ServeHTTP(w, r)
		})
	}, nil
}

func parseRules(c *setup.Controller) ([]*corsRule, error) {
	rules := []*corsRule{}

	for c.Next() {
		rule := &corsRule{Path: "/", Conf: cors.Default()}
		args := c.RemainingArgs()

		anyOrigins := false
		switch len(args) {
		case 0:
		case 2:
			rule.Conf.AllowedOrigins = strings.Split(c.Val(), ",")
			anyOrigins = true
			fallthrough
		case 1:
			rule.Path = args[0]
		default:
			return nil, c.Errf(`Too many arguments`, c.Val())
		}
		for c.NextBlock() {
			switch c.Val() {
			case "origin":
				if !anyOrigins {
					rule.Conf.AllowedOrigins = []string{}
				}
				args := c.RemainingArgs()
				for _, domain := range args {
					rule.Conf.AllowedOrigins = append(rule.Conf.AllowedOrigins, strings.Split(domain, ",")...)
				}
				anyOrigins = true
			case "methods":
				if arg, err := singleArg(c, "methods"); err != nil {
					return nil, err
				} else {
					rule.Conf.AllowedMethods = arg
				}
			case "allowCredentials":
				if arg, err := singleArg(c, "allowCredentials"); err != nil {
					return nil, err
				} else {
					var b bool
					if arg == "true" {
						b = true
					} else if arg != "false" {
						return nil, c.Errf("allowCredentials must be true or false.")
					}
					rule.Conf.AllowCredentials = &b
				}
			case "maxAge":
				if arg, err := singleArg(c, "maxAge"); err != nil {
					return nil, err
				} else {
					i, err := strconv.Atoi(arg)
					if err != nil {
						return nil, c.Err("maxAge must be valid int")
					}
					rule.Conf.MaxAge = i
				}
			case "allowedHeaders":
				if arg, err := singleArg(c, "allowedHeaders"); err != nil {
					return nil, err
				} else {
					rule.Conf.AllowedHeaders = arg
				}
			case "exposedHeaders":
				if arg, err := singleArg(c, "exposedHeaders"); err != nil {
					return nil, err
				} else {
					rule.Conf.ExposedHeaders = arg
				}
			default:
				return nil, c.Errf("Unknown cors config item: %s", c.Val())
			}
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func singleArg(c *setup.Controller, desc string) (string, error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return "", c.Errf("%s expects exactly one argument", desc)
	}
	return args[0], nil
}
