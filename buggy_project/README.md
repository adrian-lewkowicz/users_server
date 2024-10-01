# buggy project

I find problem with goroutines sync - every method call wg.Wait() and waiting for every synchronized goroutines ends.
So on the production we can have problem because we can receive more request. And the first goroutine need wait to last goroutine ends work. 
Goroutines are ended when server make all work and waiting for the next requests. In this situation we can have 10000 goroutines when 9999 goroutines end work and waiting for the last.

http.HandleFunc automatically run every request in new goroutine and manage of them. We dont' nedd use syng group in this situatuion. 