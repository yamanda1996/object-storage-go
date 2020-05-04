#!/bin/bash

for((i=0;i<10000;i++))
do
	echo "hello"$i >> ./test.log
done


