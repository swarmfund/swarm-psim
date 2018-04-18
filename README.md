[./psim/cmd/psim/main.go](./psim/cmd/psim/main.go) is the enter point,
it imports packages with all Services without directly using them, so all the services are built into the binary.

Services packages in their init() functions (see https://golang.org/doc/effective_go.html#init for details)
call the [app.RegisterService()](./psim/app/main.go) function.

Then looking the config, which file name is provided in command-line arguments ([config.yaml](./config.yaml) by default),
app decides which services to start (`services` block in the config)
and only the services present in the `services` block will be run,
usually only 1 service is run per PSIM start.

Config of each Service is normally parsed inside the setUp function of the Service,
which is normally in the main.go file in a service's package.

[figure](https://gitlab.com/distributed_lab/figure) if used to parse services' configs from *.yaml files into go structures.

## Contribution

### Logging

- Psim uses [logan](https://gitlab.com/distributed_lab/logan) for logging.

- Each error should be logged only once - on the top level of the stack.

- Each intermediate function between the error appearing and logging just only wraps
(see [errors.Wrap](https://gitlab.com/distributed_lab/logan/v3/errors/errors.go)) the error,
adding [Fields](https://gitlab.com/distributed_lab/logan/v3/fields.go) if needed.

- Only the values, which are not accessible from the upper functions along the stack,
should be attached in function to the returned error as fields.

- Field names should be named using snake_case.

- For the complex types (structs) being put into fields
[GetLoganFields()](https://gitlab.com/distributed_lab/logan/v3/fields/main.go) interface method should be implemented.


### To add a new Service

- Add side-effect import of your Service package into the [./psim/cmd/psim/main.go](./psim/cmd/psim/main.go)

- Add constant with Service name into [./psim/conf/services.go](./psim/conf/services.go)

- Add the Service name into the `services` list in the [config.yaml](./config.yaml) (must match the name just defined in conf)

- Add the specific section for the Service in the [config.yaml](./config.yaml) (section name must also match)
You normally will parse the content of this section from config into the Config go struct in your service
inside the setUp function during the `init()`.

- Be careful with parsing config structures with embedded complex structures -
[figure](https://gitlab.com/distributed_lab/figure) needs adding hook for this.

##### Make sure no sensitive of specific real data is mentioned in the [config.yaml](./config.yaml) (such as keys, secrets, real urls, etc)
