.PHONY: run

run:
	rm -f *.o search test
	shtkc search.shtk -c
	g++ search.shtk.o -o search -lm
	./search

compile: 
	rm -f *.o search test
	shtkc search.shtk -c
	g++ search.shtk.o -o search -lm

time:
	time ./search