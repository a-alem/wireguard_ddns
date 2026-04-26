# 📦 DDNS Updater
`ddns_updater` is a service that resolve a public IP and updates a DNS record with it. Meant to be deployed as a systemd one shot timer service. Written in Golang.

## 🗄️ Structure
The service is design to be modular and extensible through interface implementation. Divided into the following directories, similar to a repository pattern:
- `/`: At root level, the `main.go` file is found along with `VERSION` file and config example template file.
- `/internal`: Meant for the service internals, core config files and type declaration, including interfaces and config file parser.
- `/providers`: DNS provider implementations, where each provider occupied a subdirectory of its own.
- `/resolvers`: Public IP resolution providers implementation. Each provider is segmented into a subdirectory.

## 🚀 Build & Usage
In order to build the project, from `/ddns_updater`:
```bash
go build ./...
```

To run the project, it expects a `yaml` configuration file to be passed, a configuration example template file can be found in `/ddns_updater/config.example.yaml`.
From `/ddns_updater`:
```bash
go run . <config.yaml>
```

## 🛠️ Extending
In order to extend the usability and add more custom providers/resolvers, follow this simple guide:
- Create the respective directory under either providers/resolvers dirs
- Usually the implementation will have three files
  - `client.go`: for the actual API client that will initiate the request on the user's behalf to the remote service, the purpose is to abstract away these implementations.
  - `provider.go | resolver.go`: think of this as the `service` implementation, it's a high level abstraction that calls the lower level `client.go`
  - `types.go`: Define client interface, types for API requests and responses for marshalling/unmarshalling etc.
- Add a config type to `/internal/types.go` for your new provider/resolver and inject that into either `ProvidersConfig` or `ResolversConfig`