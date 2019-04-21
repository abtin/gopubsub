A quick sample code to create a google pubsub topic and subscription and publish/receive form/to it.

Setup:
- Create a GCP project
- From the API and services add pubsub
- Create a service account and give it pubsub role, save the credentials 

Execution:
- To create topic and subscription:
    ```
    ./gopubsub -create -projectId=<projectId> -topic=<topic> -subscription=<subscription>
    ```
- To receive messages from a topic/subscription:
    ```
    ./gopubsub -projectId=<projectId> -topic=<topic> -subscription=<subscription> -publisher=false
    ```
- To publish messages, one word at a time, on a topic/subscription:
    ```
    ./gopubsub -projectId=<projectId> -topic=<topic> -subscription=<subscription> -publisher
    ```
