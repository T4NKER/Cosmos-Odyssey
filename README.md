Peamine probleem seisnes selles, kuidas läheneda sellele, kuidas programm arvutab võimalikud marsruudid kahe punkti vahel.
Üks võimalus on leida kohe kõik võimalikud teekonnad ning need seejärel andmebaasi või mällu sisestada. See ei ole halb idee, kuid
probleem tekib siis, kui lähte -ja sihtpunkti vahele tuleb veel punkte. See tekitab tohutu andmevoo, mis 
andmebaasi läheb ning suurtemate äppidega see ei sobiks. Lisaks sellele võib ta sisestada kaks samasugust teed. Seetõttu valisin ma ka teise võimaluse. Iga päringuga antud api endpointi sisestatakse andmed otse andmebaasi ning vastavalt vajadusele need arvutatakse marsruutideks vastavalt päringule.

I used POSTMAN to test the endpoints.

15/06/2024
There is a major problem that is related to the pathfinding algorithm. This lies in the problem of many different possible paths when taking into account the fact that different legs can be provided by different companies. This combined with many database queries to get all the possible path provider combinations has resulted in the fact that the request does not get resolved successfully.
17/06/2024
This has been ameliorated by the fact that some flight times are not compatible and thus reduces the number of flights as well as preferring a single company or not. I resorted to keeping the flight times in memory as querying database for flight times took way too much time and was really memory heavy. I am not aware how it is done in real-life scenarios but probably using caching or similar techniques. For this project keeping the flight times in memory removes the requirement of queries as well as hastens the algorithm runtime. 

As sqlite does not support arrays in the database, I have decided to insert the route array as JSON, which is allowed.

Should be validation that the reservation is correct and there's no funny business