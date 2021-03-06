package main

import (
	"crypto/tls"
	"crypto/x509"
)

func main() {
	// Connecting with a custom root-certificate set.

	const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIKNTCCBh2gAwIBAgIQH3gpo9qx+XrF5yJnthZmxjANBgkqhkiG9w0BAQ0FADA0
MQswCQYDVQQGEwJERTEUMBIGA1UECgwLQ29tIFNjaWVuY2UxDzANBgNVBAMMBlN1
YiBDQTAeFw0xODExMDQxOTU1MTFaFw0xOTExMDQxOTU1MTFaMEExCzAJBgNVBAYT
AkRFMRQwEgYDVQQKDAtDb20gU2NpZW5jZTEcMBoGA1UEAwwTbWFpbC5jb20tc2Np
ZW5jZS5kZTCCBCIwDQYJKoZIhvcNAQEBBQADggQPADCCBAoCggQBAMu+JcyzaAmV
7R+yMzJ51tJUVWq/Pr6rpcyP5eZyZ6TWQcfKfoEp2LLDcnLDQRr4Mf73yKxwj9C3
FE4sS008yI0bbz7Ap3qeg3zRWBi9VYz3hO5+exdTTd91UQOwtZx5WTVoZ0+arBLy
+gX7izqQkJlixdvSD+UnF7co2oz9KB1zKPGArmWM1khC0ZvNpp1tNHI2w1b4M2Gx
Op6eIHE7BoUJYe3vXTtM9sifAOmFUlM+7TqlI9+eT3nU18D4miv7aSCIHcR2nCjK
N0OiRI9fo9Q1VzhMnUYXMj0m/WFdot9dBsYgDRbCYlVAVR2rC9oVf/bM+GMIZ1jb
bEiPoJDIKYQ47glplfc9yLv4tdynEROX4M5mp6JX+7yXkaKwhBMTxSJX6Ykfmy+o
Gqgdo3G0rx7/+Nrjmrl8hfR9l8E8Ohb48HaX1Z/1phaefMmJjywxO+esa0vAwjBG
QyP0skU4RHligUl4RNS8tD488uz7SW1H3jQfBH9ERT8vvO/7SdqvQL05yGm6GeOH
3xcKClGkCCKLxtz6Q1rVYToh9lqA7KMY8jXTahO028zVEb+9qqPH3aGkVlY1vz/T
TWIy3lyWD+sfR2l6pFAyAz4aXHnMaICZOYxJz42MZWcTli2izuK6Kag3cct4mObT
MspY5xqOAZRsSDqPohiLlzZ76O+o2z88UgXgdf9ceQaeLjZPxsh29MKlsgWgxiJf
njSI1eEXQFq6Q/Hdnj9JlafyYzLchGB7legwOOJ4SE8jcgwpzuaYMx98lfM1spzQ
sTZheqSL2r+aDWnQ0YoatFLzzR6lwMNeA7+m9UtvVaChrJ+5MCBUvBok3JDlgL1c
vFy4JhjR8l7TNg21iwJz+Gnf7dmzIIKrHg4yDIBHODikjgTf06SY9+nMg2yjljbU
0UUJii/UuRyP/j0lhM80nYfZzGyCqi6uxiMtGocuwtcVJW0lppEUrTpEAZczuyxr
fV/e0DVZY6kfVOLteM0OubnJuFl/xE8iEaLbzjo9kOs59QpNcoKftMdVfFD8G3bK
jtxNBmfA5WGMLtxep/OfO60iC73NkyeLmWGrG/4Iq+Kc4f6neybZlxuat/pxtjto
/Ovvl0MSwfklYM7jjtIZrxDJqqm9E0Wo2QzuPVNQEogj8B2bXQNwCjtGNQJHku5v
DEDD3Fu6kbcOXqzMMyZNds6/Cg39lXKCA2ihZHQMV5LujP1PwTHGiC2PY4BWbDGn
pr/ciS2xxL/+o3XMhryKXsoNMmRQK8T0MiiYzUHzE7PCESWTBqqlLvnFAb8H/D9I
gZF/ay5yMplwVu2j8uLoYUrF2chUDbenVb9DeoJvuaiu8tXpLDucL92HZGM1XhVs
lkxqLt4GGWkCAwEAAaOCATQwggEwMHcGCCsGAQUFBwEBBGswaTAzBggrBgEFBQcw
AoYnaHR0cDovL3N1Yi1jYS5jb20tc2NpZW5jZS5kZS9zdWItY2EuY3J0MDIGCCsG
AQUFBzABhiZodHRwOi8vb2NzcC5zdWItY2EuY29tLXNjaWVuY2UuZGU6OTA4MTAf
BgNVHSMEGDAWgBR+sJoqZTrrO4gZVu9kV8A5IJ1GUjAMBgNVHRMBAf8EAjAAMDgG
A1UdHwQxMC8wLaAroCmGJ2h0dHA6Ly9zdWItY2EuY29tLXNjaWVuY2UuZGUvc3Vi
LWNhLmNybDAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDgYDVR0PAQH/
BAQDAgWgMB0GA1UdDgQWBBSoIJctsNV+YILD7S6dZ/ok2nAcDDANBgkqhkiG9w0B
AQ0FAAOCBAEAEtePOnuRH74S3mpfWfURW/8YeXVWWpQ2BjGJiDNMYkXkN173Dyo9
1h6gN8ZAU/THZJ2MfjxQoO2K9qlCoxVz49XkO0T4I3Z/uvh5kxIpxLhQ47WIKTLq
Fy8Us+rU9huWSpbcuVdfGgtbV/7C5hFVddLRQct4hnS/4E/j7Y49VXiVfA0rAv/X
YcMAduDp83+HEyWxmWRTBrKvuysuAHR5hsMT/AXbw9zpQ3JUhPEEr49Sdog/xylO
cbzgwhEnx4V97j47HZcWAt63LpN/z5w0lxXeSlkS7X8S1kGco28CNxhfpzUZj85b
5ukROS61YSA8OQRqZ9erYyV2zL3S1zeVC4vveQQmMfY/S3W0tWPzDmywTs0nkCt9
NNzPSYxDsVkX6cUiChhnngn+F1bThZIonMlggJO+JA1G3F5dYF3nijePcHpbZTzp
d98ILj66oBTFI4lgMECzE91v27WhjKWu4PciDN+FZNqk91K4H+o26Hl2/8c7Mb0I
VXtthyKSoBZlA1TfsBW8gBdtIWR0QKg+6Sz5ad7NAouy0UWMOAiYzIglR/paE9G6
IWB+K23CJ57JGBIkbSr+KoMr8vk9RX1xGlS0EoQLhCsPpd+4dSlZFhbgGl31SsmD
LFpUUqU40eDdyTPFWi5z2p5BxRcdYgZCx7GN/9FkKJK05asfnkAoOVKv0mnxa57r
sCOTEwvPw9vsgql5n+BLIhTSqH6A3gX8UyuCGYytUJXOJWLKf5PyPsyEPuFQdQEi
oDx8XhY89NkrRtoKuaDUuKk1iYQlFZwNhf0kexzEeX3fhN9SUhsf81iDIe0Priym
Y3OYgafpJ9/F2pDhJBnwEHV5KQ7efQHkan9LMIiwbFVUAvx/TvANlo1BndXKr8Lt
zMczg216lJYpqcuCbkZxME4OhX+g38OfX8h6hiRTGcwbfLbPKbnPMGAYo5omkSXD
ntJISrN7emqeW2mVkh1UXCfHCfjbqFFBbrbse40MmxsgXSYgSfSJ1asnpHy7tULU
JwP//nRMvtI7d7OISgqI+82Mngfvs9m6qtrOxYYCLC1N11VXiSFb/J0IrucmJKD+
EuKfrjy5Sd+yVnPiwRx7Sv4VNdqlIrsvMOlx3RqsClcl7FRxp8/cKgt6S/BLX30Y
kxeeUtyMyDPhOah2qo6AaKBKhjg3W6ct2e2Yyd5nLmatMbZETJHDM3C2vCPZB5mA
3mOaioOW+wQW2OUjo2w8ALpDGNtNSbP8mE9RNhQJgT5H6T0T3GyLVJpFUsBO437p
4ws27TmqwiV+CIVR5TDwCJ9C21yH0DHEQYMLc/RHvTMYxbsx19Q1wnGuJKsLlBe4
IaolmWoKDhz+bGWN9bUWIK/F9fQkdqWbBQ==
-----END CERTIFICATE-----
`

	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	conn, err := tls.Dial("tcp", "mail.com-science.de:993", &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	conn.Close()
}
