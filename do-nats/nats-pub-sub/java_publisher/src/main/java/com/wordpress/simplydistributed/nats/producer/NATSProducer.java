package com.wordpress.simplydistributed.nats.producer;

import io.nats.client.Connection;
import io.nats.client.Nats;
import io.nats.client.Options;
import java.util.Date;
import java.util.concurrent.TimeUnit;

public class NATSProducer {
    
    final static String SUBJECT = "foo";
    
    public static void main(String[] args) throws Exception {
        Options.Builder builder = new Options.Builder();
        Options options = builder.timeout(3, TimeUnit.SECONDS)
                .reconnectWait(5, TimeUnit.SECONDS)
                .maxReconnect(5)
                .build();
                
        String natsServer = System.getenv().getOrDefault("NATS_SERVER", "nats://localhost:4222");
        Connection nc = Nats.connect(natsServer, options);
        
        System.out.println("Publisher connected to NATS server with ID - "+ nc.getConnectedServerId());
        
        while(true){
            String msg = "Hello World "+ new Date();
            nc.publish(SUBJECT, msg.getBytes());
            System.out.println("Sent message '" + msg + "' to NATS server with ID " + nc.getConnectedServerId());
            Thread.sleep(5000);
        }
        
    }
}
