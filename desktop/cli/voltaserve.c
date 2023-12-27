#include <stdio.h>
#include <string.h>
#ifdef _WIN32
#include <Windows.h>
#else
#include <unistd.h>
#endif

char *get_file_list(const char *path)
{
    // TODO: implement
    const char *json = "[\"/example/path\", \"/example/path/hello\", \"/example/path/hello/world\"]";
    size_t length = strlen(json) + 1;
    char *result = malloc(length);
    strcpy_s(result, length, json);
    return result;
}

char *upload_file(char *path)
{
    // TODO: implement
#ifdef _WIN32
    Sleep(3000);
#else
    sleep(3000);
#endif
    const char *json = "{\"result\": \"success\"}";
    size_t length = strlen(json) + 1;
    char *result = malloc(length);
    strcpy_s(result, length, json);
    return result;
}

int main(int argc, char *argv[])
{
    if (argc != 2)
    {
        return 1;
    }
    if (strcmp(argv[1], "--get-file-list") == 0)
    {
        char *json = get_file_list("/example");
        printf(json);
        free(json);
    }
    else if (strcmp(argv[1], "--upload-file") == 0)
    {
        char *json = upload_file("/example/file.txt");
        printf(json);
        free(json);
    }
    return 0;
}