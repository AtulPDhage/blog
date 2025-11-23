import amqp from 'amqplib';

let channel: amqp.Channel;

export const connectRabbitMQ = async () => {
    try {
        const connection = await amqp.connect({
            protocol: 'amqp',
            hostname: 'localhost',
            port: 5672,
            username: 'admin',
            password: 'admin123',
        });

        channel = await connection.createChannel();

        console.log('✅ Connected to RabbitMQ');
    } catch (error) {
        console.error('❌ Failed to connect to RabbitMQ:', error);
    }       
};

export const publishToQueue = async (queueName: string, message: any) => {
    if(!channel) {
        console.error('❌ RabbitMQ channel is not established');
        return;
    }

    await channel.assertQueue(queueName, { durable: true });

    channel.sendToQueue(queueName, Buffer.from(JSON.stringify(message)), { persistent: true });
    console.log(`📤 Message sent to queue ${queueName}`);
};

export const invalidateCacheJob = async (cacheKey: string[]) => {
    try {
        const message = {
            action : "invalidateCache",
            keys : cacheKey
        };
        await publishToQueue('cache-invalidation', message);

        console.log('✅ Cache invalidation job published to RabbitMQ');
    }
    catch (error) {
        console.error('❌ Failed to publish cache invalidation job:', error);
    }
};

