// main.go
// Accepts argument and flag input then utilizes the framework
// for building SAS Viya Docker images.
//
// Copyright 2018 SAS Institute Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
)

const usage = `
    SAS Container Recipes
    Framework to build SAS Viya Docker images and create deployments using Kubernetes.

    Single Container Arguments 
    ------------------------

      Required:

      --zip <value>
          Path to the SAS_Viya_deployment_data.zip file from your Software Order Email (SOE).
          If you do not know if your organization has a SAS license then contact
          https://www.sas.com/en_us/software/how-to-buy.html
          
            example: /path/to/SAS_Viya_deployment_data.zip

      Optional:

      --addons [<value> <value>]
          A space or comma separated list of addon names. Each requires additional configuration:
          See the GitHub Wiki: https://github.com/sassoftware/sas-container-recipes/wiki/Appendix:-Under-the-Hood
            Access Engines: access-greenplum, access-hadoop, access-odbc, access-oracle, access-pcfiles, access-postgres, access-redshift, access-teradata
            Authentication: auth-sssd, auth-demo
            Other: ide-jupyter-python3

      --base-image <value>
          The Docker image and tag from which the SAS images will build on top of
            Default: centos:latest

      --tag <value>
          Override the default tag formatted as "19.0.4-2019-03-18-09-49-38"
                                            ( <recipe-version> - <date> - <time> )

      -m|--mirror-url <value>
          The location of the mirror URL. See the Mirror Manager guide at
          https://support.sas.com/en/documentation/install-center/viya/deployment-tools/34/mirror-manager.html


    Multi-Container Arguments
    ------------------------

      Required:

      --type <value>
          Choose one of the following deployments.
            multiple: SAS Viya Programming Multi-Container deployment with Kubernetes
    	    full    : SAS Visuals based deployment with Kubernetes.

<<<<<<< HEAD
          Note: the default deployment type is 'single'.
=======
        Note: the default deployment type is 'single'.
>>>>>>> 28eb1e6acc08f5fbb7eec42f062be78e08e4525a

      --zip <value>
          Path to the SAS_Viya_deployment_data.zip file from your Software Order Email (SOE).
          If you do not know if your organization has a SAS license then contact
          https://www.sas.com/en_us/software/how-to-buy.html
          
<<<<<<< HEAD
          example: /path/to/SAS_Viya_deployment_data.zip
=======
            example: /path/to/SAS_Viya_deployment_data.zip
>>>>>>> 28eb1e6acc08f5fbb7eec42f062be78e08e4525a

      --docker-namespace <value>
          The namespace in the Docker registry where Docker
          images will be pushed to. Used to prevent collisions.

<<<<<<< HEAD
          example: mynamespace
=======
            example: mynamespace
>>>>>>> 28eb1e6acc08f5fbb7eec42f062be78e08e4525a

      --docker-registry-url <value>
          URL of the Docker registry where Docker images will be pushed to.

<<<<<<< HEAD
          example: 10.12.13.14:5000 or my-registry.docker.com
=======
            example: 10.12.13.14:5000 or my-registry.docker.com
>>>>>>> 28eb1e6acc08f5fbb7eec42f062be78e08e4525a


      Optional:

      --virtual-host 
          The Kubernetes Ingress path that defines the location of the HTTP endpoint.
          For more details on Ingress see the official Kubernetes documentation at
          https://kubernetes.io/docs/concepts/services-networking/ingress/
          
            example: user-myproject.mycluster.com

      --addons [<value> <value>]
          A space or comma separated list of addon names. Each requires additional configuration:
          See the GitHub Wiki: https://github.com/sassoftware/sas-container-recipes/wiki/Appendix:-Under-the-Hood
            Access Engines: access-greenplum, access-hadoop, access-odbc, access-oracle, access-pcfiles, access-postgres, access-redshift, access-teradata
            Authentication: auth-sssd, auth-demo
            Other: ide-jupyter-python3
<<<<<<< HEAD
        
      --base-image <value>
          The Docker image and tag from which the SAS images will build on top of.
          Default: centos:latest
=======

      --base-image <value>
          The Docker image and tag from which the SAS images will build on top of.
            Default: centos:latest
>>>>>>> 28eb1e6acc08f5fbb7eec42f062be78e08e4525a

      --mirror-url <value>
          The location of the mirror URL.See the Mirror Manager guide at
          https://support.sas.com/en/documentation/install-center/viya/deployment-tools/34/mirror-manager.html

      --tag <value>
          Override the default tag formatted as "19.0.4-2019-03-18-09-49-38"
                                            ( <recipe-version> - <date> - <time> )
    
      --workers <integer>
          Specify the number of CPU cores to allocate for the build process.
          default: Utilize all cores on the build machine

      --verbose
          Output the result of each Docker layer creation.
          default: false

      --build-only "<container-name> <container-name> ..."
          Build specific containers by providing a comma or space separated list of container names in quotes.

		  WARNING: This is meant for developers that require specific small components to rapidly be build.

            example: --build-only "consul" or --build-only "consul httpproxy"

      --version
          Print the SAS Container Recipes version and exit.


    Need some more help?

        Learn more about this project
            https://github.com/sassoftware/sas-container-recipes/
        General questions, features, and bug reports from the community
            https://github.com/sassoftware/sas-container-recipes/issues
        For FAQs, troubleshooting, and tips
            https://github.com/sassoftware/sas-container-recipes/wiki
        License Assistance
            https://support.sas.com/en/technical-support/license-assistance.html 
        License Purchases
            https://www.sas.com/en_us/software/how-to-buy.html
        Software Trials
            https://www.sas.com/en_us/trials.html 
`

func main() {
	order, err := NewSoftwareOrder()
	if err != nil {
		log.Fatal(err)
	}

	err = order.Build()
	if err != nil {
		log.Fatal(err)
	}
	order.ShowBuildSummary()
}
