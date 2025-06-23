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
                ]) {
                    script {
                        writeFile file: '.env', text: """
                        DB_USER=${env.DB_USER}
                        DB_PASS=${env.DB_PASS}
                        DB_PORT=${env.DB_PORT}
                        DB_HOST=${env.DB_HOST}
                        DB_NAME=${env.DB_NAME}
                        ELASTICHOST=${env.ELASTICHOST}
                        JWT_SECRET=${env.JWT_SECRET}
                        """
                    }
                }
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
