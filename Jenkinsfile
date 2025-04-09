#!groovy
// -*- coding: utf-8; mode: Groovy; -*-
@Library('common-utils') _

properties([
    buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '7', daysToKeepStr: '', numToKeepStr: '7')),
    disableConcurrentBuilds(),
    disableResume(),
    pipelineTriggers([
        issueCommentTrigger('[Jj]enkins.*')
    ])
]);

// Parameters

lint_version = "1.26.5"
runSastScan = false
runIacScan = false
runSecretsScan = false

def alertmanagerRegex = (~/(.*\/alertmanager\/values\/.*\.yaml)/)


slackUtils.sendBuildStatusNotification(channelId: 'G01MW7UQGQH') {
    node ('lmru-dockerhost') {
        cleanUp()
        if (runSastScan){
        sastScan() 
        }
        if (runIacScan){
        iacScan() 
        }
        if (runSecretsScan){
        secretsScan() 
        }
        stage('Checkout sources') {
            gitHubUtils.checkout()
        }
        yammlint()
        if (gitUtils.hadAnyChangesInPrFiles(filesRegex: alertmanagerRegex)) {
            amtool()
        }
        cleanUp()
    }
}

def yammlint(){
    stage('Validate YAML') {
        sh "docker run --rm -v ${workspace}:/data docker-devops.art.lmru.tech/yamllint/yamllint:${lint_version} -s ."
    }
}

def amtool(){
    stage('Validate Alertmanager config') {
        sh """
            sed -Ei '2,/config:/d' ${workspace}/common/observability/lm/alertmanager/values/*/*/config.yaml
            docker run --rm --entrypoint=/bin/ash -v ${workspace}/common/observability/lm/alertmanager/values:/tmp docker.art.lmru.tech/prometheus/alertmanager:latest -c 'amtool check-config /tmp/prod/*/config.yaml /tmp/stage/*/config.yaml'
        """
    }
}

def sastScan() {
    stage('SAST scan') {
         checkmarxUtils.runCheckmarxScan_ZeroLicense(getCXReport: true, scanFromGit: true)
    }
}

def iacScan() {
    stage('IaC scan') {
        securityUtils.iacScan(reportFormat: 'html')
    }
}

def secretsScan() {
    stage ('Secrets scan'){
        securityUtils.secretsScan(days: '14')
    }
}

def cleanUp() {
    stage ('Clean Workspace') {
        cleanWs();
    }
}
