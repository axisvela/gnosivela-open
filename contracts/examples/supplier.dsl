ontology procurement.supplier @version 1.2

entity Supplier identifiedBy canonicalSupplierId
relation suppliesTo: Supplier -> LegalEntity [many]
property riskScore: Decimal [0..1]
property status: SupplierStatus

assertion HighRiskSupplier when
  riskScore >= 0.75 and context.region in ["SG", "EU"]
requires evidence RiskAssessment
authority RiskOffice
validFor 90d
