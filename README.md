## Gasket
Gasket is an example interface for the [Cayley](https://github.com/cayleygraph/cayley) rdf graph store : 
---
A gasket is a seal which fills the space between two or more mating surfaces, generally to prevent leakage.

The gasket project is an example API written in golang featuring : 
* 'Normalized' JSON object model mapped to an rdf data store (see schema below for details)
* Integration of [Cayley](https://github.com/cayleygraph/cayley) as a go library
* Example of capturing quad meta-data using cayley v0.6
* Examples of miscellaneous plumbing required for an API in golang (external configuration, logging, unit tests for handlers)

### Items not yet implemented
This project will build, compile, and run as a functional API; however, there are certain crucial elements which are not implemented.  Those include :

* Validation
* Security
* Documentation (ala Swagger)

### Installation

#### Requirements
* mongoDb [download](https://www.mongodb.com/download-center#community)
* glide [Glide - Github](https://github.com/Masterminds/glide)

#### Config
* The project will look for a configuration path from an Environment variable at ACES_CFG.  There are two separate files that can be set.  log.toml and config.toml.  Example files for each can be found within the project within the config folder.  
* If a database server name is not found in config.toml, the project will not run

#### Build
* run 'go build' from the project directory

#### Run
* run ./gasket to start the server

### Schema
The project utilizes three JSON objects.  Node, Relation, and Metadata.  

#### Node 
Node : Is a basic object with an extensible list of properties. The JSON schema for a Node is defined as :   
```
"node": {
	"type":"object",
	"properties": {
		"id": {"type":"string"},
		"name":{"type":"string"},
		"label":{"type":"string"},
	},
	"patternProperties": {
        	"^(/[^/]+)+$": { "type": "object" }
		/**
		 * Accept any property without a '/' character
		 */
   		}
	"required":["name"]
}
```
Each node will be stored as several quads within the graph store.  The basic quad representation is 
```
<id> <schema:name> "nameParam" "label"
<id> <patternPropertyA> "propertyAParam" "label"
<id> <patternPropertyB> "propertyBParam" "label"
... etc 
```

#### Relation
Relation : will relate two nodes within a typical rdf quad.  The schema is : 
```

"relation": { 
	"type":"object",
	"properties": {
		"id": {"type":"string"},
		"sourceId":{"type":"string"},
		"type":{"type":"string"},
		"targetId":{"type":"string"},
		"label":{"type":"string"},
	}
	"required":["sourceId", "type", "targetId", "label"]
}
	
```
The relation will be stored as a quad with 
```
quad.Quad{
	Subject : sourceId,
	Predicate : type,
	Object : targetId,
	Label : label,
}
```
A JSON representation of this quad will be stored and assigned an id.  This id will be associated with the JSON object.

#### Metadata
Metadata : Provides a way of storing additional details about the relation.  The schema is similar to a node :
```
"metadata": {
	"type":"object",
	"properties": {
		"id": {"type":"string"},
		"relationId":{"type":"string"},
	},
	"patternProperties": {
        	"^(/[^/]+)+$": { "type": "object" }
		/**
		 * Accept any property without a '/' character
		 */
   		}
	"required":["relationId"]

}
```
And the method of mapping metadata to quads is similar to that used for nodes.


### License

Released under the MIT license.  See LICENSE file for more details.

