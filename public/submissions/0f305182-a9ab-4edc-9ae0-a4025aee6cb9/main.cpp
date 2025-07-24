#include<stdio.h>
#include<string.h>
#include <ctype.h>

int main() {
    char str[1001];
    scanf("%s", str);
    if(strlen(str) == 1){
    	printf("me\n");
    	return 0;
	}

    int len = strlen(str);
    for (int i = len - 1; i >= 0; i--) {
        char ch = str[i];
        if (islower(ch)) ch = toupper(ch);
        putchar(ch);
    }
    putchar('\n');
    return 0;
}
