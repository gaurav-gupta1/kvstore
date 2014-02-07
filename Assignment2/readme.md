Sir very sorry for submitting the assignment late, but me and Prabhat are trying to solve a problem of packet dropping or either blocking code by Send function. The problem is that if a server is down then trying to send it a message will block the code, or if we use DONTWAIT it drops a lot of packets. We tried many things but after all when nothing works we choose to use DONTWAIT.

Now the program works as follow:

a) I use PUSH-PULL socket to bind and listen because of the comments you give in feedback about DEALER-DEALER.

b) All the information regarding a server is stored in a structure. Information like its Peers, Sockets to those peers, its ID etc.

c) A send and a recieve function which are called in a go routine, and they run forever.

d) So simply everyone is sending messages to everyone else. And for test case I count those messages.


