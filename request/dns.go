package request 

import (
    "fmt"
	"time"
    "slices"
    "github.com/miekg/dns"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func DNSCheck(cfg db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (db.Result, error) {
    InfoLog.Printf("Running DNS check for %v \n", cfg.Name)
    var res db.Result
    
    m := new(dns.Msg)
    switch cfg.RecordType {
    case "A":
        m.SetQuestion(dns.Fqdn(cfg.Target), dns.TypeA)
    case "AAAA":
        m.SetQuestion(dns.Fqdn(cfg.Target), dns.TypeAAAA)
    case "CNAME":
        m.SetQuestion(dns.Fqdn(cfg.Target), dns.TypeCNAME)
    case "MX":
        m.SetQuestion(dns.Fqdn(cfg.Target), dns.TypeMX)
    case "SRV":
        m.SetQuestion(dns.Fqdn(cfg.Target), dns.TypeSRV)
    case "TXT":
        m.SetQuestion(dns.Fqdn(cfg.Target), dns.TypeTXT)
    default:
        res.Text = "Record Type invalid: " + cfg.RecordType
        res.Status = false
        WarningLog.Printf("Record Type invalid: %v \n", cfg.RecordType)
        return res, nil
    }

    c := new(dns.Client)
    c.Timeout = 5 * time.Second

    r, _, err := c.Exchange(m, cfg.DNSServer) // or your DNS server
    if err != nil {
        res.Text = "DNS Request failed: " + err.Error()
        res.Status = false
        WarningLog.Printf("DNS Request failed: %v \n", err.Error())
        return res, nil
    }

    if r.Rcode != dns.RcodeSuccess {
        res.Text = "DNS query failed with code: " + r.Rcode
        res.Status = false
        WarningLog.Printf("DNS query failed with code %v: %v \n", r.Rcode, err.Error())
        return res, nil
    }

    var resp []string
    switch cfg.RecordType {
    case "A":
        for _, ans := range r.Answer {
            if t, ok := ans.(*dns.A); ok {
                resp = append(resp, t.A.String())
            }
        }
    case "AAAA":
        for _, ans := range r.Answer {
            if t, ok := ans.(*dns.AAAA); ok {
                resp = append(resp, t.AAAA.String())
            }
        }
    case "CNAME":
        for _, ans := range r.Answer {
            if t, ok := ans.(*dns.CNAME); ok {
                resp = append(txts, t.Target)
            }
        }
    case "MX":
        for _, ans := range r.Answer {
            if t, ok := ans.(*dns.MX); ok {
                resp = append(resp, t.Mx)
            }
        }
    case "SRV":
        for _, ans := range r.Answer {
            if t, ok := ans.(*dns.SRV); ok {
                resp = append(resp, t.Target)
            }
        }
    case "TXT":
        for _, ans := range r.Answer {
            if t, ok := ans.(*dns.TXT); ok {
                resp = append(resp, t.Txt...)
            }
        }
    default:
        return nil, fmt.Errorf("An unexpected error occured during DNS Check.")
    }
    
    if cfg.ExpectIP != "" && !slices.Contains(resp, cfg.ExpectIP) {
        res.Text = "Error in Name resolution: Name was not resolved to the expected IP " + cfg.ExpectIP
        res.Status = false
        WarningLog.Printf("Error in Name resolution: Name was not resolved to the expected IP %v \n", cfg.ExpectIP)
        return res, nil
    }
    
    res.Status = true
    res.Text = fmt.Print("DNS resolution was successful")
    return res, nil
}
