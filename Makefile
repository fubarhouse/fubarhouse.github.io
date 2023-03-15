#!/usr/bin/make -f

default: run;

run:
	docker-compose up --build -d;

up:
	docker-compose up --build;

stop:
	docker-compose down;