def NAME
def version

pipeline {
    agent { label 'jenkins-4'}
    options {
        buildDiscarder(
            logRotator(
                numToKeepStr: '5', 
                artifactNumToKeepStr: '5'
            )
        )
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        disableConcurrentBuilds()
        skipDefaultCheckout()
    }

    stages {
        stage("Checkout") {
            steps {
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: 'refs/heads/'+env.BRANCH_NAME]],
                    extensions: [[$class: 'CloneOption', noTags: false, shallow: false, depth: 0, reference: '']],
                    userRemoteConfigs: scm.userRemoteConfigs,
                ])
                sh "git checkout ${env.BRANCH_NAME}"
                sh "git reset --hard origin/${env.BRANCH_NAME}" 
                }
        }

        stage ('Getter information') {
            steps {
                script {
                    sh 'printenv | sort'
                    NAME    = ("logstash-exporter").toLowerCase()
                    TAG = sh(returnStdout: true, script: "git tag --sort version:refname | tail -1").trim()
                    sh "echo 'tag: ${TAG}'"
                    if ( env.TAG_NAME ) {
                        version = env.TAG_NAME.toLowerCase()
                    } else if ( env.BRANCH_NAME.startsWith('PR') ) {
                        version = "${TAG}.${env.BRANCH_NAME.replace('PR-', '')}-${env.BUILD_ID}"
                    } else {
                        version = "${TAG}-${env.BRANCH_NAME}-${env.BUILD_ID}"
                    }
                    sh "echo ${version}"
                }
            }
        }

//        stage('Run Tests') {
//            agent {
//                docker { 
//                    image 'golang:1.13.0-buster'
//                    reuseNode true
//                }
//            }
//            steps {
//                sh 'go test'
//            }
//        }

        stage('Build go package') {
            agent {
                docker { 
                    image 'golang:1.13.0-buster'
                    reuseNode true
                }
            }
            steps {
                sh "go build -ldflags \"-s -X github.com/productsupcom/logstash_exporter/vendor/github.com/prometheus/common/version.Version=${version}\""
            }
        }

        stage ('Build deb package') {
            steps {
                script {
                    sh "version=${version} nfpm pkg --target ${NAME}_${version}_all.deb"
                    sh "dpkg-deb -I ${NAME}_${version}_all.deb"
                    sh "dpkg -c ${NAME}_${version}_all.deb"
                }
            }
        }

        stage ('Publish deb package') {
            when { 
                buildingTag()
             }
            steps {
                sshagent (credentials: ['jenkins-ssh']) {
                    script {
                        sh "scp -o StrictHostKeyChecking=no ${NAME}_${version}_all.deb root@aptly.productsup.com:/tmp/"
                        sh "ssh -o StrictHostKeyChecking=no root@aptly.productsup.com \"aptly repo add stable /tmp/${NAME}_${version}_all.deb && \
                           aptly publish update -passphrase-file='/root/.aptly/passphrase' -batch stable s3:aptly-productsup:debian && rm /tmp/${NAME}_${version}_all.deb\""
                    }
                }
            }
        }
    }

    post ('Cleanup') {
        cleanup {
            cleanWs deleteDirs: true
        }
    }
}

