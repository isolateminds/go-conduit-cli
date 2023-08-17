package conduit

type DockerCompose struct {
	Version  string                 `yaml:"version"`
	Services map[string]Service     `yaml:"services"`
	Volumes  map[string]interface{} `yaml:"volumes"`
	Networks map[string]interface{} `yaml:"networks"`
}

type Service struct {
	Image         string                 `yaml:"image"`
	ContainerName string                 `yaml:"container_name,omitempty"`
	Ports         []string               `yaml:"ports,omitempty"`
	Environment   map[string]string      `yaml:"environment,omitempty"`
	Networks      map[string]interface{} `yaml:"networks,omitempty"`
	Volumes       []string               `yaml:"volumes,omitempty"`
	DependsOn     []string               `yaml:"depends_on,omitempty"`
	Profiles      []string               `yaml:"profiles,omitempty"`
}
