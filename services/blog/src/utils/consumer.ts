import amqp from 'amqplib';
import { redisClient } from '../server.js';
import { sql } from './db.js';

interface CacheInvalidationMessage {
    action: string;
    keys: string[];
}

export const startCacheConsumer = async () => {
    try {
        const connection = await amqp.connect({
            protocol: 'amqp',
            hostname: 'localhost',
            port: 5672,
            username: 'admin',
            password: 'admin123',
        });

        const channel = await connection.createChannel(); 
        const queueName = 'cache-invalidation';  
        await channel.assertQueue(queueName, { durable: true });

        console.log('✅ Blog Cache Consumer connected to RabbitMQ, waiting for messages...');
        channel.consume(queueName, async (msg) => {
            if (msg) {
                try{
                    const content = JSON.parse(msg.content.toString()) as CacheInvalidationMessage;
                    console.log(`📥Blog Service Received cache invalidation message:`, content);

                    if(content.action === "invalidateCache" && content.keys.length > 0) {
                        for(const pattern of content.keys) {
                            const keys = await redisClient.keys(pattern);

                            if(keys.length > 0) {
                                await redisClient.del(keys);
                                console.log(`🗑️ Blog service Invalidated cache for pattern: ${pattern}, keys deleted: ${keys.length}`);

                                const searchQuery = "";
                                const category = "";
                                const cacheKey =  `blogs:${searchQuery}:${category}`;

                                const blogs  =  await sql`SELECT * FROM blogs ORDER BY created_at DESC`;
                                await redisClient.set(cacheKey, JSON.stringify(blogs), { EX: 3600 });
                                console.log(`🔄️ Blog Service Refreshed cache for key: ${cacheKey}`);
                            }
                        }
                    }
                    channel.ack(msg);
                } catch (error) {
                    console.error('❌ Error processing cache invalidation in blog service :', error);
                    channel.nack(msg, false, true);
                }
            }
        });
    } catch (error) {   
        console.error('❌ Cache Consumer failed to connect to RabbitMQ:', error);
    }    
};