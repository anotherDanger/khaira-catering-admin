pipeline {
    agent any
    stages {
        stage('Checkout'){
            steps{
                checkout scm
            }
        }
        
        stage('Use Credentials'){
            steps{
                withCredentials([
                    string(credentialsId: 'DB_USER', variable: 'DB_USER'),
                    string(credentialsId: 'DB_PASS', variable: 'DB_PASS'),
                    string(credentialsId: 'DB_PORT', variable: 'DB_PORT'),
                    string(credentialsId: 'DB_HOST', variable: 'DB_HOST'),
                    string(credentialsId: 'DB_NAME', variable: 'DB_NAME'),
                    string(credentialsId: 'ELASTICHOST', variable: 'ELASTICHOST'),
                    string(credentialsId: 'JWT_SECRET', variable: 'JWT_SECRET'),
                    string(credentialsId: 'GHCR_TOKEN', variable: 'GHCR_TOKEN')
                ]) {
                    script {
                        sh 'echo $GHCR_TOKEN | docker login ghcr.io -u anotherDanger --password-stdin'
                        
                        writeFile file: '.env', text: """
                        DB_USER=${DB_USER}
                        DB_PASS=${DB_PASS}
                        DB_PORT=${DB_PORT}
                        DB_HOST=${DB_HOST}
                        DB_NAME=${DB_NAME}
                        ELASTICHOST=${ELASTICHOST}
                        JWT_SECRET=${JWT_SECRET}
                        """
                    }
                }
            }
        }
        stage('Clean docker'){
            steps{
                sh 'docker compose down'
            }
        }
        stage('Start Services'){
            steps {
                sh 'docker compose up -d'
            }
        }

    }
}
