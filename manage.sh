#!/usr/bin/env bash
echo "Welcome to develop manage console for TrustArea!"
while true; do
    echo "# ------------------------ #"
    echo "Available commands:"
    echo "[0] Rebuild development docker-compose"
    echo "[1] Rebuild production docker-compose"
    echo "[Exit] Exit"
    echo "# ------------------------ #"
    echo "Enter your command: "
	read command
	case $command in
	"Exit" | "exit" )
		exit 0
		;;
	"0" )
		docker-compose -f dev.yml down --rmi all
		docker-compose -f dev.yml up -d
		;;
	"1" )
		docker-compose -f prod.yml down --rmi all
		docker-compose -f prod.yml up -d
		;;
	esac
done
