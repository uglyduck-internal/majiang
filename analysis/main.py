# get data from postgresql
import psycopg2

conn = psycopg2.connect(
    host="192.168.1.144",
    port="5432",
    database="majiang",
    user="postgres",
    password="pgpassword")

cur = conn.cursor()
cur.execute("SELECT * FROM public.fourfriends where date = '2024617';")
data = cur.fetchall()
conn.close()
print(len(data))