import ReactMarkdown from "react-markdown";
function ChatStart(props){

    let messages = props.message.slice(1,-1)

    return(
        <div className="chat chat-start">
            <div className="chat-bubble chat-bubble-success text-white ml-4 mt-4"
            >
                <ReactMarkdown>{messages}</ReactMarkdown>
                
            </div>
        </div>

    )
}

export default ChatStart;