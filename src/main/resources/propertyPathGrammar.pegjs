Expression
  = head:Term tail:(_ "/" _ Term)* {
      const xs = tail.reduce(function(acc, element) {
        return acc.concat([element[3]]);
      }, [head]);
      if (xs.length === 1) {
        return xs[0]
      } else {
        return {and: xs}
      }
    }

Term
  = head:Factor tail:(_ "|" _ Factor)* {
      const xs = tail.reduce(function(acc, element) {
        return acc.concat([element[3]]);
      }, [head]);
      if (xs.length === 1) {
        return xs[0];
      } else {
        return {or: xs}
      }
    }

Factor
  = "(" _ expr:Expression _ ")" { return expr; }
  / Iri

Iri "iri"
  = ns:[a-zA-Z0-9\-_]+ "." prop:[a-zA-Z0-9\-_]+ _ mod:["^","*"]? {
   var iri = {"iri": ns.join("") + ":" + prop.join(""), "inverse": false, "transitive": false};
   if (mod && mod === "^") {
     iri.inverse = true;
   } else if (mod && mod === "*") {
     iri.transitive = true;
   }
   return iri;
 }

_ "whitespace"
  = [ \t\n\r]*