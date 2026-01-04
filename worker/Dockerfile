FROM python:3.11-slim
RUN pip install redis lxml
COPY worker.py /
CMD ["python", "/worker.py"]