"use client";

import React, { useState, useEffect, useRef } from "react";
import { Send, Search, MessageSquare } from "lucide-react";
import { useWebSocket } from "@/contexts/WebSocketContext";

interface Message {
  id: string;
  sender_id: string;
  content: string;
  created_at: string;
}

export default function MessagesPage() {
  const { socket, isConnected } = useWebSocket();
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    if (!socket) return;

    const handleMessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "chat_message") {
          setMessages((prev) => [...prev, data.payload]);
        }
      } catch (err) {
        console.error("Failed to parse chat message", err);
      }
    };

    socket.addEventListener('message', handleMessage);

    return () => {
      socket.removeEventListener('message', handleMessage);
    };
  }, [socket]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    // To send a message, we actually hit the API. 
    // The API broadcasts it over WebSocket to the receivers.
    fetch('/api/chat/messages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      // Mock sending to a hardcoded user for the sake of UI test
      // In reality, this would use the active chat's recipient ID
      body: JSON.stringify({
        content: newMessage,
        receiver_id: "user_to_send_to" 
      })
    }).then(res => {
      if (res.ok) setNewMessage("");
    });
  };

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden flex h-[calc(100vh-140px)]">
      {/* Sidebar: Chats List */}
      <div className="w-1/3 border-r border-gray-100 flex flex-col bg-gray-50">
        <div className="p-4 border-b border-gray-100">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input 
              type="text" 
              placeholder="Search messages..." 
              className="w-full bg-white border border-gray-200 rounded-full py-2 pl-10 pr-4 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
        </div>
        <div className="overflow-y-auto flex-1 p-2 space-y-1">
          {[1, 2, 3].map(i => (
            <div key={i} className="flex items-center gap-3 p-3 rounded-lg hover:bg-white cursor-pointer transition-colors">
              <div className="w-12 h-12 rounded-full bg-indigo-100 shrink-0" />
              <div className="flex-1 min-w-0">
                <div className="flex justify-between items-baseline mb-1">
                  <h3 className="font-semibold text-gray-900 truncate">User {i}</h3>
                  <span className="text-xs text-gray-400">12:34 PM</span>
                </div>
                <p className="text-sm text-gray-500 truncate">Hey, how are you doing?</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col bg-white">
        {/* Chat Header */}
        <div className="h-16 border-b border-gray-100 flex items-center px-6 justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-indigo-100" />
            <div>
              <h2 className="font-bold text-gray-900">User 1</h2>
              <span className="text-xs text-green-500 flex items-center gap-1">
                <div className="w-1.5 h-1.5 rounded-full bg-green-500" />
                {isConnected ? "Online" : "Connecting..."}
              </span>
            </div>
          </div>
        </div>

        {/* Chat Messages */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4">
          <div className="text-center text-xs text-gray-400 my-4">Today</div>
          
          {/* Mock Received */}
          <div className="flex justify-start">
            <div className="bg-gray-100 text-gray-800 rounded-2xl rounded-tl-sm px-4 py-2 max-w-[70%]">
              Hello! This is a mock received message.
            </div>
          </div>
          
          {/* Dynamic Messages */}
          {messages.map((m, idx) => (
            <div key={idx} className="flex justify-end">
              <div className="bg-indigo-600 text-white rounded-2xl rounded-tr-sm px-4 py-2 max-w-[70%] shadow-sm">
                {m.content}
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>

        {/* Message Input */}
        <div className="p-4 bg-white border-t border-gray-100">
          <form onSubmit={handleSendMessage} className="flex gap-2">
            <input 
              type="text" 
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Type a message..." 
              className="flex-1 bg-gray-50 border border-gray-200 rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <button 
              type="submit"
              disabled={!newMessage.trim()}
              className="bg-indigo-600 text-white p-2 w-10 h-10 rounded-full flex items-center justify-center hover:bg-indigo-700 transition-colors disabled:opacity-50"
            >
              <Send className="w-4 h-4 ml-1" />
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
