#!/bin/bash

cd readlink -f $(dirname $(readlink $0))
workspacePath=`readlink -f ./../../../../..`
currentPath=`pwd`

declare -a filesToLink=(\
	'setup_environment.sh' \
	'create_links.sh'
);

echo $workspacePath
echo $currentPath

for file in ${filesToLink[@]}
do
	ln -s $currentPath/$file $workspacePath/$file
done
