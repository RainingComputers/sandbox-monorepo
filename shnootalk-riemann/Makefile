.PHONY: run

run:
	rm -f *.o search test
	shtkc search.shtk -release
	g++ search.shtk.o -o search -lm
	./search

compile: 
	rm -f *.o search test
	shtkc search.shtk -release
	g++ search.shtk.o -o search -lm

time:
	time ./search

clean:
	rm -f search
	rm -f search.shtk.o