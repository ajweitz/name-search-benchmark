from locust import HttpUser, task, between
from tasks import search_task as t

class QuickstartUser(HttpUser):
    wait_time = between(1, 5)

    @task
    def search_task(self):
        t(self.client,"/mysql/get-words-no-index")

