#yum install java -y
#yum install zlib.i686 -y

#!/bin/bash



IS_INSTALLED=$(rpm -qa |grep java)

if [ $? -eq 0 ]
then
    echo 'installed'
else
    echo 'not installed'
fi


IS_INSTALLED=$(rpm -qa |grep zlib.i686)

if [ $? -eq 0 ]
then
    echo 'installed'
else
    echo 'not installed'
fi