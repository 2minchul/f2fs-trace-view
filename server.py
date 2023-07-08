import asyncio
import json
import random

from aiohttp import web
from aiohttp.web_request import Request

app = web.Application()
routes = web.RouteTableDef()


@routes.get('/')
async def index(request: Request):
    with open('index.html') as f:
        _html = f.read()
    return web.Response(text=_html, content_type='text/html')


@routes.get('/zone/{zone_number}')
async def stream_data(request: Request):
    response = web.StreamResponse()
    response.enable_chunked_encoding()
    response.headers['Access-Control-Allow-Origin'] = '*'
    response.headers['Access-Control-Allow-Methods'] = 'GET'
    try:
        await response.prepare(request)

        for _ in range(10000):
            y = random.randint(0, 1023)
            row = [random.randint(0, 1) for _ in range(512)]
            await response.write(json.dumps([y, row]).encode('utf-8') + b'\n')
            await asyncio.sleep(0.001)

        await response.write_eof()
    except ConnectionResetError:
        pass
    return response


app.add_routes(routes)

if __name__ == '__main__':
    web.run_app(app, port=8080)
