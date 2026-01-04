FROM python:3.11-slim
RUN pip install flask redis lxml
COPY . /app
WORKDIR /app
CMD ["python", "app.py"]