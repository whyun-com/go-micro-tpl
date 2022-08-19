Set-PSDebug -Trace 1
$KAFKA_VERSION="3.2.1"
$KAFKA_URL="https://downloads.apache.org/kafka/${KAFKA_VERSION}/kafka_2.12-${KAFKA_VERSION}.tgz"
Invoke-WebRequest -uri $KAFKA_URL -outfile kafka.tgz
7z.exe x .\kafka.tgz
$kafkaDir="c:\"
7z.exe x .\kafka.tar -o"${kafkaDir}"
$base="${kafkaDir}kafka_2.12-${KAFKA_VERSION}"
Start-Process -NoNewWindow $base\bin\windows\zookeeper-server-start.bat $base\config\zookeeper.properties
Start-Sleep -Seconds 7
Start-Process -NoNewWindow $base\bin\windows\kafka-server-start.bat $base\config\server.properties

$REDIS_URL = "https://github.com/ServiceStack/redis-windows/raw/master/downloads/redis-latest.zip"

Invoke-WebRequest -uri $REDIS_URL -outfile redis-latest.zip
7z.exe x .\redis-latest.zip -oredis
Start-Process -NoNewWindow  redis\redis-server.exe redis\redis.windows.conf
npm install -g wait-on
wait-on tcp:6379
wait-on tcp:9092
netstat -ano | findstr 9092