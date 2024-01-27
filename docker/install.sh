#yum install java -y
#yum install zlib.i686 -y

#!/bin/bash



IS_INSTALLED=$(rpm -qa |grep java)

if [ $? -eq 0 ]
then
    echo 'installed'
else
    yum install java -y
fi


IS_INSTALLED=$(rpm -qa |grep zlib.i686)

if [ $? -eq 0 ]
then
    echo 'installed'
else
    yum install zlib.i686 -y
fi

