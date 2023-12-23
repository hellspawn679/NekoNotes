from flask import Flask, request
from flask_restful import Resource, Api
from flask_cors import CORS  # Import CORS
import requests
import socket
headers = {"Authorization": "Bearer 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"}
# creating the flask app 
app = Flask(__name__)
# creating an API object 
api = Api(app)
CORS(app, resources={r"/": {"origins": ""}})
# making a class for a particular resource 
# the get, post methods correspond to get and post requests 
# they are automatically mapped by flask_restful. 
# other methods include put, delete, etc. 
''' 
"inputs": "content"
'''


class T5_small(Resource):
    @staticmethod
    def post():
        API_URL = "https://api-inference.huggingface.co/models/facebook/bart-large-cnn"
        payload = request.get_json()
        response = requests.post(API_URL, headers=headers, json=payload)
        return response.json()


''' 
"inputs": {
    "question": "What is my name?",
    "context": "My name is Clara and I live in Berkeley."
  }

'''


class qna(Resource):
    @staticmethod
    def post():
        API_URL = "https://api-inference.huggingface.co/models/deepset/roberta-base-squad2"
        payload = request.get_json()
        response = requests.post(API_URL, headers=headers, json=payload)
        return response.json()


''' 
{
    "inputs": "My name is Sarah Jessica Parker but you can call me Jessica", "parameters": {"src_lang": "en_XX", "tgt_lang": "fr_XX"}
}
'''


class translate(Resource):
    @staticmethod
    def post():
        API_URL = "https://api-inference.huggingface.co/models/facebook/mbart-large-50-many-to-many-mmt"
        payload = request.get_json()
        response = requests.post(API_URL, headers=headers, json=payload)
        return response.json()

'''
{"inputs": "The answer to the universe is"}
'''
class text_generation(Resource):
    @staticmethod
    def post():
        API_URL = "https://api-inference.huggingface.co/models/gpt2"
        payload = request.get_json()
        response = requests.post(API_URL, headers=headers, json=payload)
        return response.json()


class say_hello(Resource):
    @staticmethod
    def get():
        return f"container ID: {socket.gethostname()}"
# adding the defined resources along with their corresponding urls
api.add_resource(T5_small, '/summarize')
api.add_resource(qna, '/qna')
api.add_resource(translate, '/translate')
api.add_resource(text_generation, '/text_generation')
api.add_resource(say_hello, '/')

# driver function 
if __name__ == '__main__':
    app.run(debug=False)
    
#docker compose up -d --build  --scale server=3 