import ReactMarkdown from "react-markdown";


function ChatStart(props){

    let messages =props.message.toString();
    messages = messages.slice(0,-1);
    messages = messages.replace(/\\text\{([^}]+)\}/g, '$1');
    messages = messages.replace(/\\times/g, '×');

    return(
        <div className="chat chat-start">
            <div className="chat-bubble chat-bubble-success text-white ml-4 mt-4"
            style={{ whiteSpace: "pre-wrap", wordBreak: "break-word" }}
            >
                <ReactMarkdown>{messages}</ReactMarkdown>
                
            </div>
        </div>

    )
}

export default ChatStart;