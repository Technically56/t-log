# T-Log: Terminal Log Analiz ve İzleme Aracı

T-Log, terminal tabanlı bir log izleme ve analiz aracıdır. Birden fazla log dosyasını eş zamanlı olarak izleyebilir, regex tabanlı kurallar ile analiz yapabilir ve CSV formatında raporlar oluşturabilir.

## Özellikler

- **Çoklu Dosya İzleme (Tailing):** Birden fazla log dosyasını tek bir ekranda canlı olarak izleme imkanı sunar.
- **Kural Tabanlı Analiz:** Regex kuralları ile tespit edilen tehditlerin analiz edilebilir.
- **Raporlama:** Analiz sonuçlarını CSV formatında dışa aktarılabilir.

## Kurulum 
Githubdan indirimek için:

```bash
git clone https://github.com/Technically56/t-log.git
cd t-log
```
### Docker ile Çalıştırma (Önerilen)

Docker compose ile çalıştırma:

```bash
docker compose run t-log
```
Komutun çalıştırıldığı terminal oturumunda t-log'u çalıştırır. Ayrıca /var/log klasörünün konteyner içine mount edilmesi sağlanır.
### Kaynak Koddan Çalıştırma

Go 1.24 veya üzeri gereklidir:

```bash
go run main.go
```

## Konfigürasyon Yapısı (`config.yaml`)

`config.yaml` dosyası, uygulamanın hangi log dosyalarını izleyeceğini ve hangi kuralları uygulayacağını belirler.

Örnek bir `config.yaml` yapısı:

```yaml
monitors:
  - name: "SSH Auth Monitor"          # İzlenecek servisin adı
    path: "/var/log/auth.log"         # İzlenecek log dosyasının tam yolu
    rules_path: "./rules/ssh_rules.yaml" # Kuralların bulunduğu dosya yolu
    source_color: "#ff00ff"           # Canlı izlemede kullanılacak renk (Hex veya isim)
```

- **name:** İleride eklenecek özellikler için saklanan isim.
- **path:** İzlenecek log dosyasının sistemdeki yolu. (Docker kullanıyorsanız `/var/log` konteyner içine mount edilmiştir, ayrıca izlenecek dosyaların ve dosya yollarının tamamen bu klasörde yer alması gerekir.)
- **rules_path:** Bu log dosyası için uygulanacak kurallar setinin yolu.
- **source_color:** Canlı akışta bu log kaynağından gelen satırların rengi.

## Kural Yapısı (`rules/*.yaml`)

Her log kaynağı için ayrı bir kural dosyası tanımlayabilirsiniz. Kurallar YAML formatındadır.

Örnek bir kural dosyası (`web_rules.yaml`):

```yaml
name: "Web Rules"  # Kural setinin adı
rules:
  - name: "Nginx 404"             # Kuralın adı
    regex: "\"(GET|POST|HEAD) .* HTTP/.*\" 404" # Eşleşecek Regex ifadesi
    level: "WARNING"              # Önem derecesi (INFO, WARNING, ERROR, CRITICAL,DEBUG)
    description: "Detects 404 errors" # Kuralın açıklaması

  - name: "SQL Injection Attempt"
    regex: "(?i)(union|select|insert|update|delete|drop|alter).*(from|into|table|database)"
    level: "CRITICAL"
    description: "URL içinde SQL enjeksiyon denemesi tespit eder"
```

- **name:** Kuralın adı.
- **description:** Kuralın açıklaması.
- **regex:** Go dili uyumlu regex ifadesi.
- **level:** Log seviyesi. Arayüzde renklendirme ve kritiklik bilgisi için kullanılır (CRITICAL=Kırmızı, ERROR=Kırmızı, WARNING=SARI, INFO=YEŞİL, DEBUG=MAVİ).

## Klavye Kısayolları

- **ESC:** Önceki menüye dön veya çıkış yap.
- **S:** Rapor sayfasında analiz sonucunu CSV olarak kaydeder.
- **Yukarı/Aşağı Ok:** Listelerde gezinme.
