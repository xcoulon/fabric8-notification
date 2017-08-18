#!/usr/bin/groovy
@Library('github.com/fabric8io/fabric8-pipeline-library@master')
def utils = new io.fabric8.Utils()
def initServiceGitHash
def releaseVersion
goTemplate{
  dockerNode{
    ws {
      checkout scm

      if (utils.isCI()){
        def v = goCI{
          githubOrganisation = 'fabric8-services'
          dockerOrganisation = 'fabric8'
          project = 'fabric8-notification'
          dockerBuildOptions = '--file Dockerfile.deploy'
          makeTarget = 'build test-unit-no-coverage-junit'
        }

        sh('mv /home/jenkins/go/src/github.com/fabric8-services/fabric8-notification/tmp/junit.xml `pwd`')
        junit 'junit.xml'

  /*
        container(name: "docker") {
          sh "docker tag docker.io/fabric8/fabric8-notification:${v} registry.devshift.net/fabric8-services/fabric8-notification:test"
          sh "docker push registry.devshift.net/fabric8-services/fabric8-notification:test"
        }
*/
      } else if (utils.isCD()){
        def v = goRelease{
          githubOrganisation = 'fabric8-services'
          dockerOrganisation = 'fabric8'
          project = 'fabric8-notification'
          dockerBuildOptions = '--file Dockerfile.deploy'
          makeTarget = 'build test-unit-no-coverage-junit'
        }
    
        initServiceGitHash = sh(script: 'git rev-parse HEAD', returnStdout: true).toString().trim()
      }
    }

    if (utils.isCD()){
      ws{
        container(name: 'go') {
          def gitRepo = 'openshiftio/saas-openshiftio'
          def flow = new io.fabric8.Fabric8Commands()
          sh 'chmod 600 /root/.ssh-git/ssh-key'
          sh 'chmod 600 /root/.ssh-git/ssh-key.pub'
          sh 'chmod 700 /root/.ssh-git'

          git "git@github.com:${gitRepo}"

          sh "git config user.email fabric8cd@gmail.com"
          sh "git config user.name fabric8-cd"

          def uid = UUID.randomUUID().toString()
          def branch = "versionUpdate${uid}"
          sh "git checkout -b ${branch}"

          sh "sed -i -r 's/- hash: .*/- hash: ${initServiceGitHash}/g' dsaas-services/f8-notification.yaml"

          def commitMsg = sh(script: 'git log --format=%B -n 1 HEAD', returnStdout: true).
          def message = """Update notification version to ${releaseVersion}
          
          ```
          ${commitMsg}
          ```
          """
          sh "git commit -a -m \"${message}\""
          sh "git push origin ${branch}"

          def prId = flow.createPullRequest(message, gitRepo, branch)
          flow.mergePR(gitRepo, prId)
        }
      }
    }
  
  }
}