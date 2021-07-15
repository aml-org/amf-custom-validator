#!groovy
pipeline {
    agent none
    stages {
//         stage('Test profiles (Java)') {
//             agent {
//                 dockerfile {
//                     filename 'Dockerfile'
//                     additionalBuildArgs  '--target CI-JAVA'
//                 }
//             }
//             steps {
//                 echo "Success" // Tests are actually run when building the agent in the Dockerfile
//             }
//         }
        stage('Test validator (Go)') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target CI-GO'
                }
            }
            steps {
                echo "Success" // Tests are actually run when building the agent in the Dockerfile
            }
        }
        stage('Test generated WASM (JS)') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target CI-JS'
                }
            }
            steps {
                echo "Success" // Tests are actually run when building the agent in the Dockerfile
            }
        }
    }
}
