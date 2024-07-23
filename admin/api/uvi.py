import uvicorn

if __name__ == "__main__":
    uvicorn.run("main:app", host='localhost', port=20001, reload=True, workers=1)
