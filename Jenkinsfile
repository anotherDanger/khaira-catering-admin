pipeline {
    agent any
    stages {
        stage('Checkout'){
            steps{
                checkout scm
            }
        }

        stage('Start Services'){
            steps {
                script{
                    sh 'docker compose up -d'
                }
            }
        }

        stage('Stop Services'){
            steps {
                sh 'docker compose down'
            }
        }
    }
}