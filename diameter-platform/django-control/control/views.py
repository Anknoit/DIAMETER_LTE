from django.http import HttpResponse
def index(request):
    return HttpResponse("Diameter Control UI (scaffold)")
