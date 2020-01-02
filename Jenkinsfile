def NAME
def TAG

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
                dir ('logstash_exporter') {
                    checkout scm
                }
            }
        }

        stage ('Getter information') {
            steps {
                script {
                    NAME    = ("logstash_exporter").toLowerCase()
                    TAG = env.TAG_NAME.toLowerCase()

                    sh 'printenv | sort'
                }
            }
        }

        stage('Build go package') {
            agent {
                docker { 
                    image 'golang:1.13.0-buster'
                    reuseNode true
                }
            }
            when {
              branch 'master'
            }
            steps {
                sh "go build -o ./build/${build_name} -ldflags \"-X github.com/productsupcom/logstash_exporter/cmd.appVersion=${version}\""

            }
        }


        stage ('Clean unneeded files') {
            steps {
                dir ('logstash_exporter') {
                    sh """
                        ( find . -type d -name ".git" \
                            && find . -name ".gitignore" \
                            && find . -name ".gitmodules" \
                            && find . -type d -name "test" \
                            && find . -type d -name "tests" ) | xargs rm -rf
                    """
                }
            }
        }

        stage ('Build deb package') {
            steps {
                script {
                    sh "cd logstash_exporter && TAG=${TAG} nfpm pkg --target ../${NAME}_${TAG}_all.deb"
                    sh "dpkg-deb -I ${NAME}_${TAG}_all.deb"
                    sh "dpkg -c ${NAME}_${TAG}_all.deb"
                }
            }
        }

        stage ('Publish deb package') {
            steps {
                sshagent (credentials: ['jenkins-ssh']) {
                    script {
                        PACKAGE_NAME=${NAME}_${TAG}_all.deb
                        scp ${PACKAGE_NAME} root@workspace.productsup.com:/tmp/
                        scp ${PACKAGE_NAME} root@aptly.productsup.com:/tmp/
                        ssh root@aptly.productsup.com "aptly repo add stable /tmp/${PACKAGE_NAME} && aptly publish update -passphrase-file='/root/.aptly/passphrase' -batch stable s3:aptly-productsup:debian && rm /tmp/${PACKAGE_NAME}"
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

def checkFolderForDiffs(path) {
    try {
        // git diff will return 1 for changes (failure) which is caught in catch, or
        // 0 meaning no changes 
        sh "git diff --quiet --exit-code HEAD~1..HEAD ${path}"
        return false
    } catch (err) {
        return true
    }
}
