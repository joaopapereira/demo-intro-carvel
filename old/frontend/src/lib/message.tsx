import {ConnectError, createPromiseClient} from "@bufbuild/connect";
import {MessageService} from "../../gen/board/v1/api_connect";
import {createConnectTransport} from "@bufbuild/connect-web";
import {AllMessagesResponse} from "../../gen/board/v1/api_pb";

const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
});

export async function allMessages(): Promise<AllMessagesResponse> {
    const client = createPromiseClient(MessageService, transport);
    console.log(transport)
    try {
        //@ts-ignore
        return client.allMessages({})
    } catch (err) {
        // We have to verify err is a ConnectError
        // before using it as one.
        if (err instanceof ConnectError) {
            err.code;    // Code.InvalidArgument
            err.message; // "[invalid_argument] sentence cannot be empty"
        }
    }
    return Promise.reject("error");
}

export async function addMessage(title: string, message: string) {
    const client = createPromiseClient(MessageService, transport);
    const res = await client.addMessage(
        {
            message: {
                title: title, message: message, timestamp: Math.floor(Date.now() / 1000).toString()
            }
        })
    console.log(res);
    return res
}

