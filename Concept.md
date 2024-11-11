# Projektkonzept

## 1. Lebenszyklus der Anwendung

Die Anwendung wird mithilfe einer durchdachten CI/CD-Pipeline und einer strukturierten Branch-Strategie entwickelt und bereitgestellt.

### Entwicklungsstrategie und Pipeline
1. **Feature-Entwicklung**: Entwickler arbeiten an neuen Funktionen auf Feature-Branches. Bei jedem Pull-Request wird die CI-Pipeline ausgelöst und prüft, ob alle Tests erfolgreich sind.
2. **Code Review und Merge**: Nach erfolgreichem Durchlauf der CI-Pipeline wird der Code manuell von einem Reviewer in den geschützten<sup>*</sup> Master-Branch gemerged, welcher stabil und produktionsbereit bleibt.
3. **Release-Branch und Veröffentlichung**: Ein Pull-Request vom Master wird in den geschützten Release-Branch gemerged, um gezielte Releases durchzuführen. Nach dem Merge in den Release-Branch wird die Release-Stage der Pipeline ausgelöst, und die Anwendung wird als Artefakt veröffentlicht.

<sub>* Aufgrund des gewählten Bezahlmodells kann ich die Branches nicht schützen oder bestimmte Regeln anwenden. In der realen Implementierung würde ich jedoch diese Schutzmaßnahmen in meiner Strategie einplanen. </sub>

## 2. Architektur der Infrastruktur

Die Infrastruktur wurde so gestaltet, dass sie eine zuverlässige, skalierbare und wartbare Umgebung bietet:

- **Runtime-Umgebung**: Die Container werden in Google Kubernetes Engine (GKE) gehostet. GKE bietet Funktionen für Proxy, virtuellen Host und Load Balancing und sorgt für die nötige Skalierbarkeit.
- **Überwachung und Monitoring**: Prometheus in Kombination mit Google Cloud Platform (GCP) und Grafana wird für das Monitoring eingesetzt, um die Infrastruktur und die Anwendung performant und stabil zu halten.
- **Infrastructure as Code (IaC)**: OpenTofu wird zur Definition der Infrastruktur genutzt und stellt eine reproduzierbare und versionskontrollierte Bereitstellung sicher.

## 3. Technologie-Stack

| Komponente                         | Zweck                                                 |
|------------------------------------|-------------------------------------------------------|
| **GitHub**                         | Quellcodeverwaltung, Hosten von Artefakten            |
| **GitHub Actions**                 | CI/CD, automatisierte Tests und Bereitstellung        |
| **Google Kubernetes Engine (GKE)** | Runtime-Umgebung, Proxy, VHost, Load-Balancer |
| **Prometheus oder GCP + Grafana**  | Überwachung und Monitoring                      |
| **OpenTofu**                       | IaC, Konfigurationsmanagement                         |
| **Helm**                           | Automatisierung von Installationen und Konfiguration  |

## 4. Entscheidungsbegründung

- **GitHub und GitHub Actions** bieten eine nahtlose Integration zwischen Quellcode und CI/CD. GitHub Actions ist leistungsfähig und flexibel genug, um automatisierte Tests und Bereitstellungen zu ermöglichen.
- **GKE (Google Kubernetes Engine)** stellt eine skalierbare Container-Umgebung bereit, die für Anwendungen in der Produktion optimiert ist und Load Balancing sowie Proxying integriert.
- **Prometheus/GCP und Grafana** erlauben ein umfassendes Monitoring und sorgen für Transparenz in der Infrastruktur und der Anwendungsleistung.
- **OpenTofu** (Terraform) ermöglicht die Verwaltung der Infrastruktur als Code (IaC) und gewährleistet durch Reproduzierbarkeit eine einfache Wartung und Skalierung.
- **Helm** vereinfacht die Verwaltung und das Installieren von Kubernetes-Umgebungen.

## 5. Infrastruktur- und Bereitstellungsprozesse

- **Hosting der Infrastruktur**: GKE hostet die Container und verwaltet Load Balancing und Proxying. Alle Ressourcen werden über OpenTofu definiert, um die Umgebung konsistent und reproduzierbar bereitzustellen.
- **Umgebungen**: Es gibt verschiedene Kubernetes-Namespaces oder GKE-Instanzen für Entwicklungs-, Staging- und Produktionsumgebungen.
- **Service-Bereitstellung**: Helm wird genutzt, um die Dienste auf Kubernetes bereitzustellen. Dadurch kann jede Umgebung unabhängig und nach denselben Prinzipien konfiguriert werden.

## 6. Branch-Strategie und CI/CD-Trigger

- **Branch-Strategie**: Neue Features werden auf separaten Feature-Branches entwickelt. Bei bestandener CI und Code-Review werden diese in den geschützten Master-Branch gemerged, der stets stabil und produktionsbereit bleibt. Ein Release-Branch wird genutzt, um spezifische Versionen gezielt freizugeben.
- **Pipeline-Trigger**: Die CI-Pipeline wird nur bei Pull Requests ausgelöst. Nach erfolgreichem Merge in Master wird die Pipeline erneut durchlaufen, ebenso wie bei Merge in den Release-Branch, bei dem zusätzlich die die Release-Stage aktiviert und das Artefakt veröffentlicht wird.

Dieses Konzept stellt sicher, dass die Anwendung stabil, skalierbar und kontinuierlich überprüft wird. Die Infrastruktur und die Prozesse erlauben eine effiziente Entwicklung, Überwachung und Verwaltung der Anwendung.
