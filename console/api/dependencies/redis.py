from redis.asyncio import Redis

from . import settings

redis_conn = Redis().from_url(url=settings.REDIS_URL)
