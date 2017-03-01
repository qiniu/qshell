#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/wait.h>
int main()
{
	int code=system("./main");
	printf("%d\n",WEXITSTATUS(code));
}
